package com.rakuten.plummy

import com.rakuten.plummy.errors.UnsupportedContentTypeException
import com.rakuten.plummy.files.MultiplexReader
import com.rakuten.plummy.files.MultiplexWriter
import com.rakuten.plummy.handlers.ExceptionHandler
import com.rakuten.plummy.params.PlummyHeaders
import com.rakuten.plummy.params.RawParams
import com.rakuten.plummy.params.decodeRawParams
import com.rakuten.plummy.renderers.RenderInput
import com.rakuten.plummy.renderers.Renderer
import com.rakuten.plummy.util.*
import io.undertow.Handlers
import io.undertow.Undertow
import io.undertow.UndertowOptions
import io.undertow.server.HttpServerExchange
import io.undertow.util.Headers

data class ServerConfiguration(
    val port: Int,
    val address: String,
    val renderers: List<Renderer>
)

class Server(private val config: ServerConfiguration) {
    private val renderers = config.renderers.map { it.name to it }.toMap()

    private val exceptionHandler: ExceptionHandler = routes()
        .withExceptions()
        .returnStatusCode<UnsupportedContentTypeException>(415)
        .returnStatusCode<IllegalArgumentException>(400)
        .returnStatusCode<Exception>(500) // Default handler

    fun start() {
        val rootHandler = exceptionHandler

        val server = Undertow.builder()
            .addHttpListener(config.port, config.address)
            .setHandler(rootHandler)
            // In HTTP/1.1, connections are persistent unless declared
            // otherwise.  Adding a "Connection: keep-alive" header to every
            // response would only add useless bytes.
            .setServerOption(UndertowOptions.ALWAYS_SET_KEEP_ALIVE, false)
            .build()

        server.start()
    }

    private fun render(exchange: HttpServerExchange) {
        val engineName = exchange.pathTemplateParams["engine"]
            ?: throw IllegalArgumentException("rendering engine name not specified")

        val renderer = renderers[engineName] ?: throw IllegalArgumentException("unknown rendering engine '$engineName'")

        // Parse Base64-encoded request params
        val requestParams = exchange.requestHeaders
            .getFirst(PlummyHeaders.params)
            ?.decodeRawParams()
            ?: RawParams.empty

        if (exchange.requestHeaders.getFirst(Headers.CONTENT_TYPE) != PlummyHeaders.ContentTypes.multiplex)
            throw UnsupportedContentTypeException("Content-Type must be ${PlummyHeaders.ContentTypes.multiplex}")

        exchange.requestReceiver.receiveFullBytes { recvExchange, bodyBytes ->
            val input = RenderInput(requestParams, MultiplexReader.parseFiles(bodyBytes))
            exceptionHandler.runBlocking(exchange) {
                val output = renderer.render(input)
                recvExchange.responseHeaders.also {
                    it[Headers.CONTENT_TYPE] = PlummyHeaders.ContentTypes.multiplex
                    it[PlummyHeaders.params] = output.params.encodeToString()
                }
                recvExchange.startBlocking()
                MultiplexWriter.writeFiles(recvExchange.outputStream, output.files)
            }
        }
    }

    private fun routes() = Handlers.routing(false)
        .get("/healthz") { exchange ->
            exchange.sendText("OK")
        }
        .post("/v1/{engine}/render", ::render)
}
