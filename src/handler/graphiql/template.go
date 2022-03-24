package graphiql

func graphiQlTemplate(
	url string,
	subscriptionUrl string,
) string {
	return `
		<html>
		<head>
			<title>GraphiQL</title>
			<link href="https://unpkg.com/graphiql/graphiql.min.css" rel="stylesheet" />
		</head>
		<body style="margin: 0;">
			<div id="graphiql" style="height: 100vh;"></div>

			<script crossorigin src="https://unpkg.com/react/umd/react.production.min.js"></script>
			<script crossorigin src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"></script>
			<script crossorigin src="https://unpkg.com/graphiql/graphiql.min.js"></script>

			<script>
				const fetcher = GraphiQL.createFetcher({ 
					url: '` + url + `',
					subscriptionUrl: '` + subscriptionUrl + `',
				});

				ReactDOM.render(
					React.createElement(GraphiQL, { fetcher: fetcher }),
					document.getElementById('graphiql'),
				);
			</script>
		</body>
		</html>
	`
}
