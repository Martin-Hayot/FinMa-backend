import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ["/login", "/signup"];

export default async function middleware(req: NextRequest) {
    const path = req.nextUrl.pathname;

    const isPublic = publicRoutes.includes(path);

    const token = req.cookies.get("access_token")?.value || "";

    if (path === "/" && token) {
        return NextResponse.redirect(new URL("/dashboard", req.nextUrl));
    }

    if (path === "/login" && token) {
        return NextResponse.redirect(new URL("/dashboard", req.nextUrl));
    }

    if (path === "/signup" && token) {
        return NextResponse.redirect(new URL("/dashboard", req.nextUrl));
    }

    if (!isPublic && !token) {
        return NextResponse.redirect(new URL("/login", req.nextUrl));
    }

    return NextResponse.next();
}

// Routes Middleware should not run on
export const config = {
    matcher: ["/((?!api|_next/static|_next/image|.*\\.png$).*)"],
};
