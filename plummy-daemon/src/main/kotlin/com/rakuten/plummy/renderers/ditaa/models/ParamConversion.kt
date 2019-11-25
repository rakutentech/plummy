package com.rakuten.plummy.renderers.ditaa.models

import org.stathissideris.ascii2image.core.ConversionOptions
import org.stathissideris.ascii2image.core.ProcessingOptions
import org.stathissideris.ascii2image.core.RenderingOptions
import java.awt.Color

fun GlobalParams.toConversionOptions() = ConversionOptions().apply {
    processing.applyTo(processingOptions)
    rendering.applyTo(renderingOptions)
    setDebug(debug)
}

private fun ProcessingParams.applyTo(target: ProcessingOptions) {
    target.setVerbose(verbose)
    target.characterEncoding = Charsets.UTF_8.name()
    target.setAllCornersAreRound(roundCorners)
    target.setPerformSeparationOfCommonEdges(commonEdgeSeparation)

    tabSize?.let {
        if (it < 0) throw IllegalArgumentException("tabSize must be 0 or higher")
        target.tabSize = tabSize
    }
}

private fun RenderingParams.applyTo(target: RenderingOptions) {
    target.imageType = when (format) {
        "png" -> RenderingOptions.ImageType.PNG
        "svg" -> RenderingOptions.ImageType.SVG
        else -> throw IllegalArgumentException("Unsupported rendering format: $format")
    }

    target.setAntialias(antialias)
    target.setDropShadows(dropShadows)
    target.isFixedSlope = fixedSlope

    scale?.let { target.scale = it.toFloat() }
    backgroundColor?.let { target.backgroundColor = backgroundColor.parseColor() }
    fontURL?.let { target.fontURL = it }

    // TODO: Custom shape definitions
}

fun String.parseColor(): Color = when (length) {
    6 -> Color(this.toInt(16))
    8 -> Color(
        this.substring(0, 2).toInt(16),
        this.substring(2, 4).toInt(16),
        this.substring(4, 6).toInt(16),
        this.substring(6, 8).toInt(16)
    )
    else -> throw IllegalArgumentException(
        "Cannot interpret \"$this\" as background colour." +
                "It needs to be a 6- or 8-digit hex number," +
                "depending on whether you have transparency or not (same as HTML)."
    )
}
