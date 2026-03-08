import { BrowserRouter, Route, Routes, Link } from 'react-router';

import Home from './pages/home';
import Dashboard from './pages/dashboard';

function App() {
	return (
		<BrowserRouter>
			<div className='min-h-dvh bg-base-300 flex flex-col'>
				<header className='navbar shadow-md px-4 md:px-8'>
					<div className='flex-1'>
						<Link
							to='/'
							className='btn btn-ghost font-bold tracking-tight md:text-lg lg:text-xl'
						>
							URL Shortener
						</Link>
					</div>
					<div className='flex-none'>
						<ul className='menu menu-horizontal px-1 gap-2'>
							<li>
								<Link to='/' className='rounded-lg'>
									Home
								</Link>
							</li>
							<li>
								<Link to='/dashboard' className='rounded-lg'>
									Dashboard
								</Link>
							</li>
						</ul>
					</div>
				</header>

				<main className='flex-1 container mx-auto py-8'>
					<Routes>
						<Route index element={<Home />} />
						<Route path='/dashboard' element={<Dashboard />} />
					</Routes>
				</main>

				<footer className='footer footer-center p-4 text-base-content'>
					<p className='inline-block'>
						<span>Copyright © 2026 - All right reserved.</span>
						<span className='mx-2'>|</span>

						<a
							href='https://vilas-gannaram.github.io/'
							target='_blank'
							rel='noopener noreferrer'
							className='underline underline-offset-4 font-semibold uppercase inline-block'
						>
							Vilas Gannaram
						</a>
					</p>
				</footer>
			</div>
		</BrowserRouter>
	);
}

export default App;
