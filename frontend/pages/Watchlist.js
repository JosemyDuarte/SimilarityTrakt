import {signOut, useSession} from "next-auth/react";
import {useRouter} from "next/router";
import styles from "../styles/Watchlist.module.css";
import {useEffect, useState} from "react";

export default function Watchlist() {
    const {data: session, status} = useSession();
    const router = useRouter();
    const [watchlist, setWatchlist] = useState([]);

    async function retrieveWatchlist() {
        try {
            const res = await fetch("http://localhost:8080/trakt/watchlist", {
                headers: {
                    Authorization: `Bearer ${session.accessToken}`,
                },
            });
            if (res.ok) {
                const data = await res.json();
                setWatchlist(data.movies);
            } else {
                console.error(`Failed to retrieve watchlist. Status code: ${res.status}`);
            }
        } catch (error) {
            console.error("Failed to retrieve watchlist.", error);
        }
    }

    // 

    useEffect(() => {
        if (session) {
            retrieveWatchlist();
        }
    }, [session]);

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
                <h1 className={styles.heading}>Your Watchlist</h1>
                {watchlist && watchlist.length > 0 ? (
                    <div>
                        <ul className={styles.movieList}>
                            {watchlist.map((item) => (
                                <li className={styles.movieListItem} key={item.id}>
                                    <h3 className={styles.movieTitle}>
                                        <span className={styles.movieTypeIcon}></span>
                                        {item.movie.title} ({item.movie.year})
                                    </h3>
                                    <p className={styles.movieInfo}>{item.type}</p>
                                </li>
                            ))}
                        </ul>
                        <button className={styles.button} onClick={() => {
                            router.push("/Similarity");
                        }}>
                            Find Similarity
                        </button>

                    </div>
                ) : (
                    <p className={styles.message}>No items in watchlist.</p>
                )}
            </main>
        </div>
    );

}
