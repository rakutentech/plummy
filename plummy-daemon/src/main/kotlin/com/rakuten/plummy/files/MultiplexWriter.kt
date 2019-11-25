package com.rakuten.plummy.files

import com.rakuten.plummy.util.ByteArraySlice
import java.io.OutputStream

object MultiplexWriter {
    fun writeFiles(stream: OutputStream, files: List<FileData>) {
        files.forEach(stream::writeFileData)
    }
}

fun OutputStream.writeChunk(chunk: ByteArraySlice) {
    write(chunk.size.littleEndianBytes())
    write(chunk.data, chunk.start, chunk.size)
}

fun OutputStream.writeChunk(chunk: ByteArray) {
    write(chunk.size.littleEndianBytes())
    write(chunk)
}

fun OutputStream.writeFileData(fileData: FileData) {
    writeChunk(fileData.name.toByteArray())
    writeChunk(fileData.metadata)
    writeChunk(fileData.contents)
}

private fun Int.littleEndianBytes(): ByteArray {
    val result = ByteArray(4)
    result[0] = (this and 0xFF).toByte()
    result[1] = ((this shr 8) and 0xFF).toByte()
    result[2] = ((this shr 16) and 0xFF).toByte()
    result[3] = ((this shr 24) and 0x7F).toByte()
    return result
}