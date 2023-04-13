import {signIn, signOut, useSession} from "next-auth/react"
import {useRouter} from 'next/router';
import styles from '../styles/Login.module.css'

export default function Login() {
    const {data: session} = useSession()
    const router = useRouter();

    if (session) {
        return (
            <div className={styles.container}>
                <h1 className="title">Dreamly</h1>
                <div className={styles.content}>
                    <h2> Signed in as {session.user.name} <br/></h2>
                    <div className={styles.btns}>
                        <button className={styles.button} onClick={() => router.push('/Profile')}>User Profile</button>
                        <button className={styles.button} onClick={() => {
                            signOut()
                        }}>Sign out
                        </button>
                    </div>
                </div>
            </div>
        )
    }
    return (
        <div className={styles.container}>
            <h1 className="title">Dreamly</h1>
            <div className={styles.content}>
                <h2> You are not signed in!</h2>
                <button className={styles.button} onClick={() => signIn('trakt', {callbackUrl: '/Profile'})}>Sign in
                </button>
            </div>
        </div>
    )
}
