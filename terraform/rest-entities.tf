variable "entities" {
  description = "Map of entities and operations."
  type        = map(any)
  default = {
    # TODO: This can likely be done way more elegantly
    "scaffold-create" = {
      entity    = "scaffold",
      operation = "create"
    },
    "scaffold-delete" = {
      entity    = "scaffold",
      operation = "delete"
    },
    "scaffold-get" = {
      entity    = "scaffold",
      operation = "get"
    },
    "scaffold-list" = {
      entity    = "scaffold",
      operation = "list"
    },
    "scaffold-update" = {
      entity    = "scaffold",
      operation = "update"
    },
    "profile-create" = {
      entity    = "profile",
      operation = "create"
    },
    "profile-delete" = {
      entity    = "profile",
      operation = "delete"
    },
    "profile-get" = {
      entity    = "profile",
      operation = "get"
    },
    "profile-list" = {
      entity    = "profile",
      operation = "list"
    },
    "profile-update" = {
      entity    = "profile",
      operation = "update"
    },
  }
}
