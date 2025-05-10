"use client";

import * as React from "react";
import {
    IconDashboard,
    IconDatabase,
    IconFileWord,
    IconReport,
    IconSearch,
    IconSettings,
    // IconUsers,
} from "@tabler/icons-react";

// import { NavDocuments } from "@/components/navigation/nav-documents";
import { NavMain } from "@/components/navigation/nav-main";
import { NavSecondary } from "@/components/navigation/nav-secondary";
import { NavUser } from "@/components/navigation/nav-user";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "@/components/ui/sidebar";
import { GalleryVerticalEnd, LucideIcon } from "lucide-react";
import Link from "next/link";
import { useUserStore } from "@/store/useUser";

const navigation = {
    navMain: [
        {
            title: "Dashboard",
            url: "/dashboard",
            icon: IconDashboard,
        },
    ],
    navSecondary: [
        {
            title: "Settings",
            url: "/settings",
            icon: IconSettings,
        },
        {
            title: "Search",
            url: "#",
            icon: IconSearch,
        },
    ],
    documents: [
        {
            name: "Data Library",
            url: "#",
            icon: IconDatabase,
        },
        {
            name: "Reports",
            url: "#",
            icon: IconReport,
        },
        {
            name: "Word Assistant",
            url: "#",
            icon: IconFileWord,
        },
    ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    const { user } = useUserStore();
    return (
        <Sidebar collapsible="offcanvas" {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton
                            asChild
                            className="data-[slot=sidebar-menu-button]:!p-1.5"
                        >
                            <Link href="/dashboard">
                                <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
                                    <GalleryVerticalEnd className="size-4" />
                                </div>
                                <span className="text-base font-semibold">
                                    FinMa
                                </span>
                            </Link>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <NavMain
                    items={
                        navigation.navMain as unknown as {
                            title: string;
                            url: string;
                            icon?: LucideIcon | undefined;
                        }[]
                    }
                />
                {/* <NavDocuments
                    items={
                        data.documents as unknown as {
                            name: string;
                            url: string;
                            icon: LucideIcon;
                        }[]
                    }
                /> */}
                <NavSecondary
                    items={
                        navigation.navSecondary as unknown as {
                            title: string;
                            url: string;
                            icon: LucideIcon;
                        }[]
                    }
                    className="mt-auto"
                />
            </SidebarContent>
            <SidebarFooter>
                <NavUser
                    user={{
                        name: user?.firstName + " " + user?.lastName,
                        email: user?.email || "",
                        avatar: user?.avatar || "",
                    }}
                />
            </SidebarFooter>
        </Sidebar>
    );
}
