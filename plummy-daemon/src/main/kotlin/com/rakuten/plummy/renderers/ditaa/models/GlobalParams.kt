package com.rakuten.plummy.renderers.ditaa.models

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration

@Serializable
data class GlobalParams(
    val rendering: RenderingParams = RenderingParams(),
    val processing: ProcessingParams = ProcessingParams(),
    val debug: Boolean = false
) {
    companion object {
        private val json = Json(JsonConfiguration.Stable)

        fun parse(raw: ByteArray) =
            if (raw.isEmpty()) GlobalParams()
            else json.parse(serializer(), raw.toString(Charsets.UTF_8))
    }
}

@Serializable
data class RenderingParams(
    val antialias: Boolean = true,
    val dropShadows: Boolean = true,
    val fixedSlope: Boolean = true,
    val format: String? = null,
    val scale: Double? = null,
    val backgroundColor: String? = null,
    val fontURL: String? = null
)

@Serializable
data class ProcessingParams(
    val verbose: Boolean = false,
    val commonEdgeSeparation: Boolean = true,
    val roundCorners: Boolean = false,
    val tabSize: Int? = null
)

