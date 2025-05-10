"use server";

import axios from "axios";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function logout() {
    const cookieStore = await cookies();
    const res = await axios.post(
        `${process.env.BACKEND_URL}/api/auth/logout`,
        {},
        { withCredentials: true }
    );
    if (res.status !== 200) {
        return { success: false };
    }
    cookieStore.delete("access_token");
    cookieStore.delete("refresh_token");
    return redirect("/login");
}
