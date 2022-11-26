variable env {
  description = "deployment environment"
  type = string
  default = "dev"
}

variable api_stage {
  description = "stage for API gateway"
  type = string
  default = "dev"
}

variable log_level {
  description = "zerolog log level"
  type = number
  default = 1 # INFO

  # Zerolog levels:
  # TRACE    : -1
  # DEBUG    :  0
  # INFO     :  1
  # WARN     :  2
  # ERROR    :  3
  # FATAL    :  4
  # PANIC    :  5
  # DISABLED :  7
}