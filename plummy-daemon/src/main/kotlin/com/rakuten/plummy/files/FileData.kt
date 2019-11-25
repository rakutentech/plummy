package com.rakuten.plummy.files

import com.rakuten.plummy.util.ByteArraySlice
import java.io.OutputStream

class FileData(val name: String, val metadata: ByteArraySlice, val contents: ByteArraySlice)

