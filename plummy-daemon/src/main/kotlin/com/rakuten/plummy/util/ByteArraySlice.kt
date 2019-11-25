package com.rakuten.plummy.util

import java.io.ByteArrayInputStream
import java.nio.ByteBuffer
import java.nio.charset.Charset

class ByteArraySlice(val data: ByteArray, val start: Int, val size: Int) {
    fun buffer(): ByteBuffer = ByteBuffer.wrap(data, start, size)
    fun stream() = ByteArrayInputStream(data, start, size)
    fun decodeString(charset: Charset = Charsets.UTF_8) = String(data, start, size, charset)

    companion object {
        fun of(data: ByteArray) = ByteArraySlice(data, 0, data.size)
        val empty = ByteArraySlice(ByteArray(0), 0, 0)
    }
}

fun ByteArray.asSlice() = slice(0, size)
fun ByteArray.slice(start: Int, size: Int = 0) = ByteArraySlice(this, start, size)

