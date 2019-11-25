package com.rakuten.plummy.params

import net.sourceforge.plantuml.FileFormat

object FileFormats {
    fun fromName(name: String): FileFormat = when (name) {
        "braille" -> FileFormat.BRAILLE_PNG
        "eps" -> FileFormat.EPS
        "eps:text" -> FileFormat.EPS_TEXT
        "latex" -> FileFormat.LATEX
        "latex:nopreamble" -> FileFormat.LATEX_NO_PREAMBLE
        "pdf" -> FileFormat.PDF
        "png" -> FileFormat.PNG
        "scxml" -> FileFormat.SCXML
        "svg" -> FileFormat.SVG
        "txt" -> FileFormat.ATXT
        "utxt" -> FileFormat.UTXT
        "xmi" -> FileFormat.XMI_STANDARD
        "xmi:argo" -> FileFormat.XMI_ARGO
        "xmi:start" -> FileFormat.XMI_STAR
        else -> throw IllegalArgumentException("Unknown file format $name")
    }
}