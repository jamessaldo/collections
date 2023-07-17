import { z } from "zod";
import {
  createTRPCRouter,
  publicProcedure,
  protectedProcedure,
} from "@/server/api/trpc";

export const greetingRouter = createTRPCRouter({
  getGreeting: publicProcedure.query(() => {
    return "Hello, world!";
  }),

  hello: protectedProcedure
    .input(z.object({ text: z.string() }))
    .query(({ input }) => {
      return `Hello, ${input.text}`;
    }),

  getSecretMessage: protectedProcedure.query(() => {
    return "you can now see this secret message!";
  }),
});
