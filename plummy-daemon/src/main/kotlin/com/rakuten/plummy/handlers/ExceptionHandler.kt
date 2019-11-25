package com.rakuten.plummy.handlers

import io.undertow.server.HttpHandler
import io.undertow.server.HttpServerExchange
import java.util.concurrent.CopyOnWriteArrayList

/**
 * Handler that dispatches to a given handler and allows mapping exceptions
 * to be handled by additional handlers.  The order the exception handlers are
 * added is important because of inheritance.  Add all child classes before their
 * parents in order to use different handlers.
 */
class ExceptionHandler(private val parentHandler: HttpHandler) : HttpHandler {
    private val exceptionHandlers = CopyOnWriteArrayList<Holder<*>>()

    override fun handleRequest(exchange: HttpServerExchange) =
        handleExceptionsIn(exchange) {
            parentHandler.handleRequest(exchange)
        }

    fun runBlocking(exchange: HttpServerExchange, block: () -> Unit) {
        if (exchange.isInIoThread) {
            // Dispatch in blocking thread
            exchange.dispatch( HttpHandler { handleExceptionsIn(it) { block() } })
        } else {
            block() // Call as-is
        }
    }

    fun <T : Exception> handle(clazz: Class<T>, handler: (HttpServerExchange, T) -> Unit): ExceptionHandler {
        exceptionHandlers.add(Holder(clazz, handler))
        return this
    }

    inline fun <reified T : Exception> handle(noinline handler: (HttpServerExchange, T) -> Unit) =
        handle(T::class.java, handler)


    private inline fun handleExceptionsIn(exchange: HttpServerExchange, block: () -> Unit) {
        try {
            block()
        } catch (e: Exception) {
            if (exceptionHandlers.none { it.tryHandle(exchange, e) })
                throw e // Rethrow if no handler was found
        }

    }

    private data class Holder<T : Exception>(val clazz: Class<T>, val handler: (HttpServerExchange, T) -> Unit) {
        fun <U : Exception> tryHandle(exchange: HttpServerExchange, e: U) =
            if (clazz.isInstance(e)) {
                @Suppress("UNCHECKED_CAST")
                (this as Holder<U>).handler.invoke(exchange, e)
                true
            } else false
    }
}
