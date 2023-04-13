import {signOut, useSession} from "next-auth/react";
import {useRouter} from "next/router";
import styles from "../styles/Watchlist.module.css";
import {useEffect, useState} from "react";

export default function Similarity() {
    const {data: session, status} = useSession();
    const router = useRouter();
    const [similarity, seSimilarity] = useState(0);
    const [username, setUsername] = useState("");

    async function findSimilarity() {
        try {
            const res = await fetch("http://localhost:8080/trakt/similarity", {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${session.accessToken}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    other_username: username,
                }),
            });
            if (res.ok) {
                const data = await res.json();
                console.log(data);
                seSimilarity(data.score);
            } else {
                console.error(`Failed to retrieve similarity. Status code: ${res.status}`);
            }
        } catch (error) {
            console.error("Failed to retrieve similarity.", error);
        }
    }

    useEffect(() => {
        // This will run every time similarity changes
        console.log('Similarity changed: ', similarity);
    }, [similarity]);


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
            <header className={styles.header}>
                <div className={styles.user}>
                    <img src={session.user.image} alt={session.user.name}/>
                    <div>
                        <p>{session.user.name}</p>
                        <div>
                            <button className={styles.button} onClick={() => signOut({callbackUrl: '/'})}>Sign Out
                            </button>
                            <button className={styles.button} onClick={() => {
                                router.push("/Profile");
                            }}> User Profile
                            </button>
                        </div>
                    </div>
                </div>
            </header>
            <main className={styles.content}>
                <h1 className={styles.heading}>Similarity</h1>
                {similarity ? (
                    // Shows smiliarity score
                    <p className={styles.message}>Similarity score: {similarity}</p>
                ) : (
                    // ask for username to find similarity
                    <div>
                        <div className={styles.message}> Enter a username to find similarity:
                            <input type="text" className={styles.input} value={username}
                                   onChange={(e) => setUsername(e.target.value)}/>
                        </div>
                        <button className={styles.button} onClick={() => {
                            findSimilarity();
                        }}>Find Similarity
                        </button>
                    </div>

                )}
            </main>
        </div>
    );

}
