import {signOut, useSession} from "next-auth/react";
import {useRouter} from "next/router";
import styles from "../styles/Profile.module.css";

export default function Profile() {
    const {data: session, status} = useSession();
    const router = useRouter();

    if (status === "loading") {
        return <p>Loading...</p>;
    }

    // If no session exists, display access denied message
    if (!(session && session.user && session.user.name)) {
        return (
            <div className={styles.container}>
                <h1 className="title">Dreamly</h1>
                <div className={styles.content}>
                    <h1>Access Denied</h1>
                    <p>You must be signed in to view this page.</p>
                    <button className={styles.button} onClick={() => {
                        router.push("/Login");
                    }}>
                        Sign In
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <h1 className="title">Dreamly</h1>
            <div className={styles.content}>
                <h1>Hello!</h1>
                <h2>You are signed in as {session.user.name} </h2>
                {<img src={session.user.image}/>}
                <button className={styles.button} onClick={() => {
                    router.push("/Login");
                }}>
                    Go back
                </button>
                <button className={styles.button} onClick={() => {
                    router.push("/Watchlist");
                }}>
                    Watchlists
                </button>
                <button className={styles.button} onClick={() => signOut({callbackUrl: '/'})}>
                    Sign out
                </button>
            </div>
        </div>
    );
}
