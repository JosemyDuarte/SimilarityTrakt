import Head from 'next/head'
import styles from '../styles/Home.module.css'
import {useRouter} from 'next/router';

export default function Home() {

    const router = useRouter();

    return (
        <div className={styles.container}>
            <Head>
                <title>Create Next App</title>
                <meta name="description" content="Generated by create next app"/>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <main className={styles.main}>
                <h1 className={styles.title}>
                    Welcome to Dreamly!
                </h1>

                <p className={styles.description}>
                    Get started by signing in{' '}
                    <code className={styles.code}>with your Tratk Account</code>

                    <button
                        className={styles.loginButton}
                        onClick={() => router.push('/Login')}> Login
                    </button>
                </p>
            </main>
        </div>
    )
}
