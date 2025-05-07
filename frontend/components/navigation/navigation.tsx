import Link from "next/link";

const Navigation = () => {
    return (
        <nav className="flex w-full fixed top-0 flex-row justify-between p-10 text-white">
            <Link href={"/"} className="text-4xl">
                FinMa
            </Link>
            <div className="flex flex-row gap-x-18 text-xl">
                <Link href={"/login"}>Login</Link>
                <Link href={"/signup"}>Sign up</Link>
            </div>
        </nav>
    );
};

export default Navigation;
