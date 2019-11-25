package com.rakuten.plummy.renderers.ditaa

import com.rakuten.plummy.files.FileData
import com.rakuten.plummy.params.RawParams
import com.rakuten.plummy.renderers.RenderInput
import com.rakuten.plummy.renderers.RenderOutput
import com.rakuten.plummy.renderers.Renderer
import com.rakuten.plummy.renderers.ditaa.models.GlobalParams
import com.rakuten.plummy.renderers.ditaa.models.toConversionOptions
import com.rakuten.plummy.util.ByteArraySlice
import com.rakuten.plummy.util.formats.InputFormat
import com.rakuten.plummy.util.formats.replaceExtension
import org.stathissideris.ascii2image.core.ConversionOptions
import org.stathissideris.ascii2image.core.ProcessingOptions
import org.stathissideris.ascii2image.core.RenderingOptions
import org.stathissideris.ascii2image.core.RenderingOptions.ImageType
import org.stathissideris.ascii2image.graphics.BitmapRenderer
import org.stathissideris.ascii2image.graphics.Diagram
import org.stathissideris.ascii2image.graphics.SVGRenderer
import org.stathissideris.ascii2image.text.TextGrid
import java.io.ByteArrayInputStream
import java.io.ByteArrayOutputStream
import java.io.PrintStream
import javax.imageio.ImageIO

class DitaaRenderer : Renderer {
    override val name: String
        get() = "ditaa"

    override fun render(input: RenderInput): RenderOutput {
        val params = GlobalParams.parse(((input.params.bytes)))
        val options = params.toConversionOptions()
        val ctx = RenderContext()
        val rendered = ctx.renderFile(input.files.first(), options)
        return RenderOutput(
            RawParams(ctx.stdoutToBytes()),
            files = listOf(rendered)
        )
    }

    private fun RenderContext.renderFile(inputFile: FileData, options: ConversionOptions): FileData {
        val grid = prepareGrid(inputFile.contents.data, options.processingOptions)
        val diagram = Diagram(grid, options)
        return diagram.render(inputFile.name, options.renderingOptions)
    }

    private fun RenderContext.prepareGrid(data: ByteArray, options: ProcessingOptions) = TextGrid().apply {
        // Add markup for custom shape
        options.customShapes?.let { addToMarkupTags(it.keys) }

        val lineBuilders = arrayListOf<StringBuilder>()
        ByteArrayInputStream(data).reader().forEachLine { line ->
            lineBuilders += StringBuilder(line)
        }
        initialiseWithLines(lineBuilders, options)

        if (options.printDebugOutput())
            stdout.println(debugString)
    }

    private fun Diagram.render(name: String, options: RenderingOptions): FileData =
        when (options.imageType!!) {
            ImageType.PNG -> {
                val image = BitmapRenderer().renderToImage(this, options)
                val output = ByteArrayOutputStream()
                if (!ImageIO.write(image, "png", output))
                    throw RuntimeException("Unexpected failure at ImageIO.write()")
                createResultData(name, ".png", output.toByteArray())
            }

            ImageType.SVG -> SVGRenderer()
                .renderToImage(this, options)
                .let { renderedText ->
                    createResultData(name, ".svg", renderedText.toByteArray())
                }
        }

    private fun createResultData(name: String, extension: String, data: ByteArray): FileData {
        return FileData(
            name = name.replaceExtension(InputFormat.Ditaa, extension),
            metadata = ByteArraySlice.empty,
            contents = ByteArraySlice.of(data)
        )
    }
}

private class RenderContext {
    private val underlyingStdoutStream = ByteArrayOutputStream()
    val stdout = PrintStream(underlyingStdoutStream)

    fun stdoutToBytes() = underlyingStdoutStream.toByteArray()
}
