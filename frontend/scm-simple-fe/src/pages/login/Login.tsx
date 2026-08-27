import { useState } from 'react'
import type { SubmitEvent } from 'react'
import { Navigate, useNavigate } from 'react-router'
import { AUTH_STORAGE_KEY } from '../../App'
import { Terminology } from '../../word_dicts/words'



export function Login() {
    const navigate = useNavigate()
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')

    if (sessionStorage.getItem(AUTH_STORAGE_KEY) === 'true') {
        return <Navigate to="/items/search" replace />
    }

    function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
        event.preventDefault()

        if (!username.trim() || !password) {
            setError('Enter your username and password to continue.')
            return
        }

        sessionStorage.setItem(AUTH_STORAGE_KEY, 'true')
        navigate('/items/search', { replace: true })
    }

    return (
        <main className="flex min-h-screen items-center justify-center bg-slate-100 px-6 py-12">
            <section className="flex-1 m-1 w-full max-w-md rounded-lg border border-slate-300 bg-white p-8 shadow-sm">
                <div className="mb-8">
                    <p className="font-mono text-md text-center uppercase tracking-[0.2em] text-slate-500">
                        {Terminology.scm_simple}
                    </p>
                    <h1 className="mt-3 text-2xl font-semibold text-slate-900">Sign in</h1>
                    <p className="mt-2 text-sm text-sl  text-slate-600">
                        Authenticate to access inventory operations.
                    </p>
                </div>

                <form className="space-y-5" onSubmit={handleSubmit}>
                    <label className="block text-sm font-medium text-slate-700">
                        Username
                        <input
                            autoComplete="username"
                            className="mt-2 block w-full rounded-md border border-slate-300 px-3 py-2.5 text-slate-900 outline-none focus:border-slate-600 focus:ring-2 focus:ring-slate-200"
                            onChange={(event) => setUsername(event.target.value)}
                            required
                            type="text"
                            value={username}
                        />
                    </label>

                    <label className="block text-sm font-medium text-slate-700">
                        Password
                        <input
                            autoComplete="current-password"
                            className="mt-2 block w-full rounded-md border border-slate-300 px-3 py-2.5 text-slate-900 outline-none focus:border-slate-600 focus:ring-2 focus:ring-slate-200"
                            onChange={(event) => setPassword(event.target.value)}
                            required
                            type="password"
                            value={password}
                        />
                    </label>

                    {error && <p className="text-sm text-red-700" role="alert">{error}</p>}

                    <button
                        className="w-full rounded-md bg-slate-900 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2"
                        type="submit"
                    >
                        Sign in
                    </button>
                </form>
            </section>
        </main>
    )
}

