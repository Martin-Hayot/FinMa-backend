import { z } from "zod";

export const SignUpSchema = z
    .object({
        firstName: z.string().min(1, { message: "Firstname is required" }),
        lastName: z.string().min(1, { message: "Lastname is required" }),
        email: z
            .string()
            .min(1, { message: "Email is required" })
            .email({ message: "Invalid email" }),
        password: z.string().min(6, {
            message: "Password must be at least 6 characters",
        }),
        confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
        message: "Passwords don't match",
        path: ["confirmPassword"],
    });

export const LoginSchema = z.object({
    email: z
        .string()
        .min(1, { message: "Email is required" })
        .email({ message: "Invalid email" }),
    password: z.string().min(1, { message: "Password is required" }),
});
