package com.rakuten.plummy.cli

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.parameters.options.default
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.types.int
import com.rakuten.plummy.Server
import com.rakuten.plummy.ServerConfiguration
import com.rakuten.plummy.renderers.ditaa.DitaaRenderer
import com.rakuten.plummy.renderers.plantuml.PlantUmlRender

class DaemonCommand : CliktCommand(name = "java -jar plummy-daemon-{version}.jar") {
    private val port: Int by option("-p", "--port", help = "Port to listen on").int().default(4545)

    override fun run() {
        hideUI()

        val config = ServerConfiguration(
            port = port,
            address = "localhost",
            renderers = listOf(PlantUmlRender(), DitaaRenderer())
        )
        Server(config).start()
    }

    private fun hideUI() {
        System.setProperty("apple.awt.UIElement", "true")
    }
}
