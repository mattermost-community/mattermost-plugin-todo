import React from 'react';

type PostPermalinkProps = {
    postPermalink: string;
}

const PostPermalink = ({postPermalink}: PostPermalinkProps) => (
    <a
        className='theme markdown_link'
        href={postPermalink}
        rel='noreferrer'
        data-link={postPermalink}
    >
        <span data-link={postPermalink}>
            {'Permalink'}
        </span>
    </a>
);

export default PostPermalink;
