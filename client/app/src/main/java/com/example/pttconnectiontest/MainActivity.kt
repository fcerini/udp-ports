package com.example.pttconnectiontest

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress


class MainActivity : AppCompatActivity() {

    private val remoteHost = "10.74.231.208"// ""172.31.21.8"//""190.2.45.173" //"119.8.74.219"
    private val pingPort = 64742 //64742
    private val pongPort = 64749 //64742

    private var pingSocket: DatagramSocket? = null
    private var pongSocket: DatagramSocket? = null

    private var resultado = ""
    private lateinit var textLog: TextView;

    private var iniciar = true

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        textLog = findViewById(R.id.textLog)
        textLog.text = "Prueba de conectividad UDP"

        val button = findViewById<Button>(R.id.button)
        button.setOnClickListener {
            textLog.text = "Ejecutando prueba...\n"
            resultado = ""

            Thread(Runnable {
                Thread.sleep(500)
                test()

                runOnUiThread {
                    textLog.text = resultado
                }

            }).start()

        }
    }


    private fun test() {
        try {

            if (iniciar){
                iniciar = false

                pingSocket = DatagramSocket()
                pingSocket!!.broadcast = true

                pongSocket =  DatagramSocket()//DatagramSocket(6666)
                pongSocket!!.broadcast = true
                pongSocket!!.connect(InetAddress.getByName(remoteHost), pongPort)

                resultado += "connect OK \n"

            }

            var ping = ByteArray(100)
            val sendPacket =
                DatagramPacket(ping, ping.size, InetAddress.getByName(remoteHost), pingPort)
            pingSocket?.send(sendPacket)
            resultado += "ping OK \n"

            //val sendPacket2 =
            //    DatagramPacket(ping, ping.size, InetAddress.getByName(remoteHost), pongPort)
            //pongSocket?.send(sendPacket2)


            val buffer = ByteArray(1024)
            val packet = DatagramPacket(buffer, buffer.size)

            for (i in 1..10) { // igual espera el FIN
                pongSocket?.receive(packet)

                val respuesta = String(buffer, Charsets.UTF_8)

                resultado += "> Server: " + respuesta + "\n"

                if (respuesta.startsWith("FIN")) {
                    break
                }
            }

        } catch (e: Exception) {
            resultado += "ERR " + e.message + "\n"
        }

    }

}