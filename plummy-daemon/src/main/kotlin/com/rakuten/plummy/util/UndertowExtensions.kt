package com.rakuten.plummy.util

import com.rakuten.plummy.handlers.ExceptionHandler
import com.rakuten.plummy.models.ErrorResponse
import io.undertow.server.HttpHandler
import io.undertow.server.HttpServerExchange
import io.undertow.util.HeaderMap
import io.undertow.util.Headers
import io.undertow.util.HttpString
import io.undertow.util.PathTemplateMatch
import kotlinx.serialization.KSerializer
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonConfiguration

typealias HandlerFunc = (HttpServerExchange) -> Unit

val HttpServerExchange.pathTemplateParams: Map<String, String>
    get() = getAttachment(PathTemplateMatch.ATTACHMENT_KEY)?.parameters ?: emptyMap()

fun HttpServerExchange.sendText(response: String) {
    responseHeaders[Headers.CONTENT_TYPE] = "text/plain; charset=UTF-8"
    responseSender.send(response)
}

fun <T> HttpServerExchange.sendJson(response: T, serializer: KSerializer<T>) {
    responseHeaders[Headers.CONTENT_TYPE] = "application/json; charset=UTF-8"
    responseSender.send(json.stringify(serializer, response))
}

operator fun HeaderMap.set(name: HttpString, value: String) {
    put(name, value)
}

// TODO: Refactor withException to use own code which works in blocking mode
fun HttpHandler.withExceptions(): ExceptionHandler = ExceptionHandler(this)

inline fun <reified T : Exception> ExceptionHandler.returnStatusCode(statusCode: Int): ExceptionHandler =
    handle<T> { exchange, e ->
        val response = ErrorResponse(
            type = e.javaClass.simpleName,
            description = e.message
        )
        exchange.statusCode = statusCode
        exchange.sendJson(response, ErrorResponse.serializer())
        e.printStackTrace()
    }

private val json = Json(JsonConfiguration.Stable)