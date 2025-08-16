import environment from "#dashboard-api/config/environment.ts";
import express from 'express'
import { pinoHttp } from "pino-http";


const { port } = environment
const app = express()

app.use(pinoHttp())

app.get('/health', (_request, response) => {
  _request.log.info("log")
  response.send("Ok nice! ")
})

app.listen(port)



