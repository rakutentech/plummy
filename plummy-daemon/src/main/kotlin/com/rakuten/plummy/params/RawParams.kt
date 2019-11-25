package com.rakuten.plummy.params

import com.rakuten.plummy.util.fromBase64
import com.rakuten.plummy.util.toBase64

data class RawParams(val bytes: ByteArray) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as RawParams

        if (!bytes.contentEquals(other.bytes)) return false

        return true
    }

    override fun hashCode(): Int {
        return bytes.contentHashCode()
    }

    companion object {
        val empty = RawParams(ByteArray(0))
    }

    fun encodeToString() = bytes.toBase64()
}

fun String.decodeRawParams() = RawParams(fromBase64())
