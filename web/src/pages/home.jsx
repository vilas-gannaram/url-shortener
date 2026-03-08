import { useState } from 'react';

function Home() {
	const [url, setUrl] = useState('');
	const [shortUrl, setShortUrl] = useState('');
	const [error, setError] = useState('');

	const [isSubmitting, setIsSubmitting] = useState(false);

	const sendRequest = async (e) => {
		e.preventDefault();
		setIsSubmitting(true);
		const response = await fetch('/api/shorten', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ url }),
		});

		if (response.ok) {
			const data = await response.json();
			if (data.error) {
				setError(data.error);
			} else {
				setShortUrl(data.short_url);
			}
		} else {
			setError('Something went wrong');
		}

		setIsSubmitting(false);
	};

	return (
		<div className='flex flex-col justify-center items-center px-4 md:px-8'>
			<div className='max-w-3xl mt-[25dvh] w-full'>
				<form onSubmit={sendRequest} className='md:min-w-3xl'>
					<label htmlFor='url'>Enter your destination URL</label>

					<div className='md:flex items-center mt-2 gap-x-2.5 space-y-2.5 md:space-y-0'>
						<div className='flex-1'>
							<input
								id='url'
								required
								type='url'
								className='input validator w-full focus:outline-none focus:ring-0'
								placeholder='https://example.com/my-long-url'
								value={url}
								onChange={(e) => setUrl(e.target.value)}
								pattern='^(https?://)?([a-zA-Z0-9]([a-zA-Z0-9-].*[a-zA-Z0-9])?.)+[a-zA-Z].*$'
								title='Must be valid URL'
							/>
							{/* <p className='validator-hint'>Must be valid URL</p> */}
						</div>

						<button disabled={isSubmitting} type='submit' className='btn'>
							{isSubmitting ? (
								<>
									<span className='loading loading-spinner loading-sm'></span>
									<span>Creating...</span>
								</>
							) : (
								<span>Create Short Link</span>
							)}
						</button>
					</div>
				</form>

				<div className='mt-4 flex flex-col gap-y-2'>
					{/* shortUrl */}
					{shortUrl && (
						<div role='alert' className='alert alert-success alert-soft'>
							<span>Short URL:</span>

							<a href={shortUrl} target='_blank' rel='noopener noreferrer'>
								{shortUrl}
							</a>
						</div>
					)}

					{/* Error message */}
					{error && (
						<div role='alert' className='alert alert-error alert-soft'>
							<span>{error}</span>
						</div>
					)}
				</div>
			</div>
		</div>
	);
}

export default Home;
