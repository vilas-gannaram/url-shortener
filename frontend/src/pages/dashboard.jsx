import { useEffect, useState } from 'react';

const Dashboard = () => {
	const [urls, setUrls] = useState([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState('');

	useEffect(() => {
		const fetchUrls = async () => {
			try {
				const response = await fetch('/api/urls');
				if (!response.ok) {
					throw new Error('Failed to fetch URLs');
				}
				const data = await response.json();
				setUrls(data || []);
			} catch (err) {
				setError(err.message);
			} finally {
				setLoading(false);
			}
		};

		fetchUrls();
	}, []);

	const copyToClipboard = (text) => {
		navigator.clipboard.writeText(text);
		// You could add a toast here if available
	};

	if (loading) {
		return (
			<div className='flex justify-center items-center h-64'>
				<span className='loading loading-spinner loading-lg text-primary'></span>
			</div>
		);
	}

	if (error) {
		return (
			<div className='p-8'>
				<div role='alert' className='alert alert-error'>
					<span>Error: {error}</span>
				</div>
			</div>
		);
	}

	const totalClicks = urls.reduce((acc, curr) => acc + curr.click_count, 0);

	return (
		<div className='p-8 max-w-6xl mx-auto'>
			<div className='lg:flex justify-between items-end mb-8 space-y-4'>
				<div>
					<h2 className='text-3xl font-bold text-base-content'>Stats</h2>
					<p className='text-base-content/60'>
						Manage and track your shortened links
					</p>
				</div>

				<div className='stats shadow bg-base-200'>
					<div className='stat'>
						<div className='stat-title'>Total Links</div>
						<div className='stat-value text-primary'>{urls.length}</div>
					</div>
					<div className='stat'>
						<div className='stat-title'>Total Clicks</div>
						<div className='stat-value text-secondary'>{totalClicks}</div>
					</div>
				</div>
			</div>

			<div className='overflow-x-auto'>
				<table className='table table-xs lg:table-md lg:table-pin-cols'>
					<thead>
						<tr>
							<th>Original URL</th>
							<th>Short URL</th>
							<th>Clicks</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						{urls.length === 0 ? (
							<tr>
								<td colSpan='4'>
									No URLs shortened yet. Go to home and create one!
								</td>
							</tr>
						) : (
							urls.map((url) => {
								const shortFull = `${import.meta.env.VITE_API_URL}/${url.short_code}`;
								return (
									<tr key={url.id}>
										<td>
											<a
												href={url.original_url}
												target='_blank'
												rel='noopener noreferrer'
												className='hover:text-primary transition-colors'
												title={url.original_url}
											>
												{url.original_url}
											</a>
										</td>
										<td>
											<span className='badge badge-ghost badge-sm'>
												/{url.short_code}
											</span>
										</td>
										<td>
											<span className='badge badge-ghost badge-sm'>
												{url.click_count}
											</span>
										</td>
										<td className=''>
											<div className='flex justify-start gap-2'>
												<button
													onClick={() => copyToClipboard(shortFull)}
													className='btn btn-ghost btn-sm tooltip'
													data-tip='Copy Link'
												>
													<svg
														xmlns='http://www.w3.org/2000/svg'
														className='h-4 w-4'
														fill='none'
														viewBox='0 0 24 24'
														stroke='currentColor'
													>
														<path
															strokeLinecap='round'
															strokeLinejoin='round'
															strokeWidth='2'
															d='M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3'
														/>
													</svg>
												</button>
												<a
													href={shortFull}
													target='_blank'
													rel='noopener noreferrer'
													className='btn btn-ghost btn-sm tooltip'
													data-tip='Open Link'
												>
													<svg
														xmlns='http://www.w3.org/2000/svg'
														className='h-4 w-4'
														fill='none'
														viewBox='0 0 24 24'
														stroke='currentColor'
													>
														<path
															strokeLinecap='round'
															strokeLinejoin='round'
															strokeWidth='2'
															d='M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14'
														/>
													</svg>
												</a>
											</div>
										</td>
									</tr>
								);
							})
						)}
					</tbody>
				</table>
			</div>
		</div>
	);
};

export default Dashboard;
