package com.rakuten.plummy.util

import java.nio.ByteBuffer

fun ByteArray.asBuffer(): ByteBuffer = ByteBuffer.wrap(this)
