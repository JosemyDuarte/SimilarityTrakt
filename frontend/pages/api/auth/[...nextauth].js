import NextAuth from "next-auth/next";
import TraktProvider from "next-auth/providers/trakt";


export default NextAuth({
    secret: process.env.AUTH_SECRET,
    providers: [
        TraktProvider({
            // Staging API // TODO change to production
            authorization: {
                url: "https://api-staging.trakt.tv/oauth/authorize",
                params: {scope: "public"}, // when default, trakt returns auth error
            },
            token: "https://api-staging.trakt.tv/oauth/token",
            clientId: process.env.TRAKT_CLIENT_ID,
            clientSecret: process.env.TRAKT_CLIENT_SECRET,
            userinfo: {
                async request(context) {
                    const res = await fetch("https://api-staging.trakt.tv/users/me?extended=full", {
                        headers: {
                            Authorization: `Bearer ${context.tokens.access_token}`,
                            "trakt-api-version": "2",
                            "trakt-api-key": context.provider.clientId,
                        },
                    })

                    if (res.ok) return await res.json()

                    throw new Error("Expected 200 OK from the userinfo endpoint")
                },
            },
            profile(profile) {
                return {
                    id: profile.ids.slug,
                    name: profile.name,
                    email: null, // trakt does not provide user emails
                    image: profile.images.avatar.full, // trakt does not allow hotlinking
                }
            },
        }),
    ],
    callbacks: {
        async jwt({token, user, account, profile, isNewUser}) {
            // Save access token to token object
            if (account?.access_token) {
                token.accessToken = account.access_token;
            }
            return token;
        },
        async session({session, token}) {
            // Add access token to session object
            session.accessToken = token.accessToken;
            return session;
        },
    },
});
