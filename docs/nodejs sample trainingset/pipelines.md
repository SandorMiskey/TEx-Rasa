# default

## onIntent(doubt.qna)

// compiler=javascript
return input.answer = input.utterance+" is a noce webpage "



## onIntent(joke.chucknorris)
// compiler=javascript
const something = request.get('http://api.icndb.com/jokes/random');
if (something && something.value && something.value.joke) {
  input.answer = something.value.joke;
}

