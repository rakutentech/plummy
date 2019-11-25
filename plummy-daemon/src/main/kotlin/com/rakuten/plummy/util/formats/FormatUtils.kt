package com.rakuten.plummy.util.formats

sealed class InputFormat(val extensions: Set<String>) {
    constructor(vararg extensions: String) : this(extensions.toSet())

    object PlantUML : InputFormat("plantuml", "puml")
    object Ditaa : InputFormat("ditaa")
}

fun String.removeExtension(inputFormat: InputFormat): String {
    val extensionIndex = lastIndexOf('.')
    if (extensionIndex == -1) return this

    val extension = substring(extensionIndex + 1)
    return if (inputFormat.extensions.contains(extension))
        this.substring(0, extensionIndex)
    else this // No matching extension found
}

fun String.replaceExtension(inputFormat: InputFormat, newExtension: String) =
    removeExtension(inputFormat) + newExtension
