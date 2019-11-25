package com.rakuten.plummy.renderers

import com.rakuten.plummy.files.FileData
import com.rakuten.plummy.params.RawParams

interface Renderer {
    val name: String
    fun render(input: RenderInput): RenderOutput
}

data class RenderInput(val params: RawParams, val files: List<FileData>)
data class RenderOutput(val params: RawParams, val files: List<FileData>)
