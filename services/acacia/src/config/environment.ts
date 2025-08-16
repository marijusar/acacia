import z from "zod"

const schema = z.object({ port: z.string() })

const unsafeEnvironment = {
  port: process.env.PORT
}

const environment = schema.parse(unsafeEnvironment)

export default environment


