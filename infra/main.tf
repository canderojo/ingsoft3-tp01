terraform {
  required_version = ">= 1.5"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0"
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
  }
}

provider "neon" {}
provider "render" {}

resource "neon_project" "preprod" {
  name = "turnero-preprod"

  # Sin esta línea el primer apply falla: el provider pide 86400 s por defecto
  # y el tope del plan gratuito es 21600.
  history_retention_seconds = 21600
}

resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "turnero_preprod_owner"
}

resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "turnero_preprod"
  owner_name = neon_role.preprod.name
}

locals {
  # La cadena de conexión se arma con lo que devolvió Neon. Nadie la copia y pega.
  conn = "postgres://${neon_role.preprod.name}:${neon_role.preprod.password}@${neon_project.preprod.database_host}/${neon_database.preprod.name}?sslmode=require"
}

resource "render_web_service" "api" {
  name   = "turnero-api-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      # Sin la etiqueta pegada: la etiqueta va en su propio campo.
      image_url = var.image_repo_api
      tag       = var.image_tag
    }
  }

  env_vars = {
    DATABASE_URL = { value = local.conn }
  }
}

resource "render_web_service" "front" {
  name   = "turnero-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  env_vars = {
    # La URL sale de la API de preprod. Se le saca la barra final porque nginx
    # descarta la ruta del pedido si el valor termina en "/".
    BACKEND_URL  = { value = trimsuffix(render_web_service.api.url, "/") }
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}
