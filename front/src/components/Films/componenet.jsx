import {reviews} from "../cinemas/assets/mock.js";

export const Cinemas = (cinema) => {
    return (
        <div>
            <h3>Films</h3>
            <div key={movie.id}>
                {movie.title}
            </div>
            ))

            <div>
                <h2>Reviews</h2>
                {reviews.map((review) => (
                    <>
                        <div>
                            <h4>{review.user.id} {review.user.name} </h4>
                        </div>
                        <div>
                            {review.text}
                            <p>Оценка фильма: {review.rating}</p>
                        </div>
                    </>
                ))}
            </div>
        </div>)

}