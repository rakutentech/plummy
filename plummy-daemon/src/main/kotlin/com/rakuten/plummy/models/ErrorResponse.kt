package com.rakuten.plummy.models

import kotlinx.serialization.Serializable

@Serializable
data class ErrorResponse(val type: String, val description: String?)
