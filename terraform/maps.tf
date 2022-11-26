variable "entities" {
  description = "Map of entities and operations."
  type        = map(any)
  default = {
    # TODO: This can likely be done way more elegantly
    "profile-create" = {
      entity    = "profile",
      operation = "create"
    },
    # "profile-delete" = {
    #   entity    = "profile",
    #   operation = "delete"
    # },
    "profile-get" = {
      entity    = "profile",
      operation = "get"
    },
    # Probably never going to want this one
    # "profile-list" = {
    #   entity    = "profile",
    #   operation = "list"
    # },
    # "profile-update" = {
    #   entity    = "profile",
    #   operation = "update"
    # },
    # etc...
    # "scaffold-create" = {
    #   entity    = "scaffold",
    #   operation = "create"
    # },
    # "scaffold-delete" = {
    #   entity    = "scaffold",
    #   operation = "delete"
    # },
    # "scaffold-get" = {
    #   entity    = "scaffold",
    #   operation = "get"
    # },
    # "scaffold-list" = {
    #   entity    = "scaffold",
    #   operation = "list"
    # },
    # "scaffold-update" = {
    #   entity    = "scaffold",
    #   operation = "update"
    # },
  }
}

variable "log_levels" {
  description = "Map of log levels to Zerolog ints"
  type = map(number)
  default = {
    "TRACE" = -1
    "DEBUG" = 0
    "INFO" = 1
    "WARN" = 2
    "ERROR" = 3
  }
}