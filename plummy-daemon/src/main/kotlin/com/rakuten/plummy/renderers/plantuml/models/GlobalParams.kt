package com.rakuten.plummy.renderers.plantuml.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration

@Serializable
data class GlobalParams(
    val format: String? = null,
    val defines: Map<String, String> = emptyMap(),
    val configLines: List<String> = emptyList()
) {
    companion object {
        private val json = Json(JsonConfiguration.Stable)

        fun parse(raw: ByteArray) =
            if (raw.isEmpty()) GlobalParams()
            else json.parse(serializer(), raw.toString(Charsets.UTF_8))
    }
}
