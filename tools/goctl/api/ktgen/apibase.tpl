package {{.}}

import com.google.gson.Gson
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter
import java.net.HttpURLConnection
import java.net.URL

const val SERVER = "http://localhost:8080"

suspend fun apiRequest(
    method: String,
    uri: String,
    body: Any = "",
    onOk: ((String) -> Unit)? = null,
    onFail: ((String) -> Unit)? = null,
    eventually: (() -> Unit)? = null
) = withContext(Dispatchers.IO) {
    val url = URL(SERVER + uri)
    with(url.openConnection() as HttpURLConnection) {
        connectTimeout = 3000
        requestMethod = method
        doInput = true
        if (method == "POST" || method == "PUT" || method == "PATCH") {
            setRequestProperty("Content-Type", "application/json; charset=utf-8")
            doOutput = true
            val data = when (body) {
                is String -> {
                    body
                }
                else -> {
                    Gson().toJson(body)
                }
            }
            val wr = OutputStreamWriter(outputStream)
            wr.write(data)
            wr.flush()
        }

         try {
            if (responseCode >= 400) {
                BufferedReader(InputStreamReader(errorStream)).use {
                    val response = it.readText()
                    onFail?.invoke(response)
                }
                return@with
            }
            //response
            BufferedReader(InputStreamReader(inputStream)).use {
                val response = it.readText()
                onOk?.invoke(response)
            }
        } catch (e: Exception) {
            e.message?.let { onFail?.invoke(it) }
        }
    }
    eventually?.invoke()
}
