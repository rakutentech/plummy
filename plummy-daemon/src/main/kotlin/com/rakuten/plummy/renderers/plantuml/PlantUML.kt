package com.rakuten.plummy.renderers.plantuml

import com.rakuten.plummy.util.ByteArraySlice
import com.rakuten.plummy.files.FileData
import com.rakuten.plummy.params.FileFormats
import com.rakuten.plummy.params.RawParams
import com.rakuten.plummy.renderers.RenderInput
import com.rakuten.plummy.renderers.RenderOutput
import com.rakuten.plummy.renderers.Renderer
import com.rakuten.plummy.renderers.plantuml.models.GlobalParams
import com.rakuten.plummy.util.asSlice
import com.rakuten.plummy.util.formats.InputFormat
import com.rakuten.plummy.util.formats.replaceExtension
import net.sourceforge.plantuml.*
import net.sourceforge.plantuml.error.PSystemError
import net.sourceforge.plantuml.preproc.Defines
import java.io.ByteArrayOutputStream
import java.io.Reader

class PlantUmlRender : Renderer {
    override val name = "plantuml"

    override fun render(input: RenderInput): RenderOutput {
        val params = GlobalParams.parse(((input.params.bytes)))
        val format = params.format
            ?.let { FileFormats.fromName(it) } // Convert file format from string
            ?: FileFormat.PNG                  // Default format

        val inputFile = input.files.firstOrNull() ?: throw IllegalArgumentException("No files")
        val blockUml = readBlockUmls(inputFile, params).first()
        val outputFile = blockUml.renderDiagram(inputFile.name, format)
        return RenderOutput(RawParams.empty, listOf(outputFile))
    }

    private fun readBlockUmls(fileData: FileData, globalParams: GlobalParams) =
        readBlockUmls(fileData.name, fileData.contents.stream().reader(), globalParams)

    private fun readBlockUmls(filename: String, reader: Reader, globalParams: GlobalParams): List<BlockUml> =
        BlockUmlBuilder(
            globalParams.configLines,
            "UTF-8",
            globalParams.defines.toDefines(),
            reader,
            FileSystem.getInstance().currentDir,
            filename
        ).blockUmls

    private fun BlockUml.renderDiagram(name: String, format: FileFormat): FileData {
        val diagram = diagram
        if (diagram is PSystemError) {
            // TODO: Return/track errors in metadata
            // we still want to return error as image too
        }

        // TODO: Log data / return log metadata
        val baos = ByteArrayOutputStream(256 * 1024 * 1024)
        diagram.exportDiagram(baos, 0, FileFormatOption(format))

        val baseFileName = name.replaceExtension(InputFormat.PlantUML, format.fileSuffix)
        return FileData(baseFileName + format.fileSuffix, ByteArraySlice.empty, baos.toByteArray().asSlice())
    }

    private fun Map<String, String>.toDefines() = Defines
        .createEmpty()!!
        .also { result ->
            forEach { (name, value) -> result.define(name, listOf(value), false, null) }
        }
}
