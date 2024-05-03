package com.example.pttconnectiontest

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress


class MainActivity : AppCompatActivity() {

    private val remoteHost = "190.2.45.173"//"10.74.231.208"// ""172.31.21.8" //"119.8.74.219"
    private val pingPort = 64748 //64742
    private val datosPort = 64747 //64742

    private var pingSocket: DatagramSocket? = null
    private var datosSocket: DatagramSocket? = null

    private var resultado = ""
    private lateinit var textLog: TextView;


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
                iniciar()
                loopPing()

                updateUI()

            }).start()

            Thread(Runnable {
                Thread.sleep(500)
                loopDatos()
                updateUI()

            }).start()

        }
    }
    private fun updateUI() {
        runOnUiThread {
            textLog.text = resultado
        }
    }
    private fun iniciar() {

        pingSocket = DatagramSocket(9000)
        pingSocket!!.broadcast = true
        pingSocket!!.connect(InetAddress.getByName(remoteHost), pingPort)

        datosSocket =  DatagramSocket(9001)//DatagramSocket(6666)
        datosSocket!!.broadcast = true
        datosSocket!!.connect(InetAddress.getByName(remoteHost), datosPort)

        resultado += "connect OK \n"
        updateUI()
        Thread.sleep(2000)

    }

    private fun loopPing() {
        try {

            val ping = ByteArray(100)
            val sendPacket =
                DatagramPacket(ping, ping.size, InetAddress.getByName(remoteHost), pingPort)
            pingSocket?.send(sendPacket)
            resultado += "ping OK \n"

            val buffer = ByteArray(1024)
            val packet = DatagramPacket(buffer, buffer.size)

            for (i in 1..10) { // igual espera el FIN
                pingSocket?.receive(packet)

                val respuesta = String(buffer, Charsets.UTF_8)

                resultado += "> Server: " + respuesta + "\n"

                updateUI()
                if (respuesta.startsWith("FIN")) {
                    break
                }
            }

        } catch (e: Exception) {
            resultado += "ERR " + e.message + "\n"
        }

    }

    private fun loopDatos() {
        try {

            var ping = ByteArray(1)

            val sendPacket2 =
                DatagramPacket(ping, ping.size, InetAddress.getByName(remoteHost), datosPort)
            datosSocket?.send(sendPacket2)

            resultado += "ping Datos OK \n"

            val buffer = ByteArray(1024)
            val packet = DatagramPacket(buffer, buffer.size)

            for (i in 1..10) { // igual espera el FIN
                datosSocket?.receive(packet)

                val respuesta = String(buffer, Charsets.UTF_8)

                resultado += "> Server: " + respuesta + "\n"

                updateUI()

                if (respuesta.startsWith("FIN")) {
                    break
                }
            }

        } catch (e: Exception) {
            resultado += "ERR " + e.message + "\n"
        }

    }

}