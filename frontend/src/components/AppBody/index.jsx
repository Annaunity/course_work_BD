import "./index.scoped.css";

function AppBody(props) {

    return (
        <div className="container">
            <main>
                {props.children}
            </main>
        </div >
    );
}

export default AppBody;
