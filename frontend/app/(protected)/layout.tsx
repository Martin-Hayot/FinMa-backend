"use client";

import { AppSidebar } from "@/components/app-sidebar";
import { SiteHeader } from "@/components/site-header";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useUserStore } from "@/store/useUser";
import axios from "axios";
import { useEffect } from "react";

const DashboardLayout = ({ children }: { children: React.ReactNode }) => {
    const setUser = useUserStore((state) => state.setUser);

    useEffect(() => {
        const fetchUser = async () => {
            const res = await axios.get("http://localhost:8080/api/auth/me", {
                withCredentials: true,
            });

            console.log(res.data.user);
            if (res.status === 200) {
                setUser(res.data.user);
            }
        };

        fetchUser();
    }, [setUser]);

    return (
        <SidebarProvider>
            <AppSidebar variant="sidebar" />
            <SidebarInset>
                <SiteHeader />
                {children}
            </SidebarInset>
        </SidebarProvider>
    );
};

export default DashboardLayout;
