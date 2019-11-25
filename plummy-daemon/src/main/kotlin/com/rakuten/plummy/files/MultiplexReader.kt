package com.rakuten.plummy.files

import com.rakuten.plummy.util.ByteArraySlice
import com.rakuten.plummy.util.toHex
import java.nio.ByteBuffer
import java.nio.ByteOrder

object MultiplexReader {
    fun parseFiles(data: ByteArray): List<FileData> {
        val reader = DataReader(data)
        val files = mutableListOf<FileData>()
        while (reader.hasMore()) {
            files += reader.readFile()
        }
        return files
    }
}

private class DataReader(val data: ByteArray) {
    private val buffer = ByteBuffer.wrap(data)
        .apply {
            order(ByteOrder.LITTLE_ENDIAN)
        }

    fun hasMore() = buffer.hasRemaining()

    private fun readSize(description: String): Int {
        val offset = buffer.position().toHex()
        if (buffer.remaining() < 4)
            throw IllegalArgumentException("Bad size for $description at byte offset 0x$offset")

        @Suppress("UsePropertyAccessSyntax")
        return buffer.getInt()
    }

    private fun readChunk(description: String): ByteArraySlice {
        val size = readSize(description)
        val data = ByteArraySlice(data, buffer.position(), size)

        // Advance buffer position
        buffer.position(buffer.position() + size)
        return data
    }

    private fun readStringChunk(description: String) =
        readChunk(description).decodeString()

    fun readFile(): FileData {
        val name = readStringChunk("filename")
        val metadata = readChunk("file metadata")
        val fileBytes = readChunk("file content")
        return FileData(name, metadata, fileBytes)
    }
}

