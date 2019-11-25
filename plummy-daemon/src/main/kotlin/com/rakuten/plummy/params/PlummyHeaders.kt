package com.rakuten.plummy.params

import io.undertow.util.HttpString

object PlummyHeaders  {
    val params = HttpString("x-plummy-params")

    object ContentTypes {
        const val multiplex = "application/x-plummy-multiplex"
    }
}

