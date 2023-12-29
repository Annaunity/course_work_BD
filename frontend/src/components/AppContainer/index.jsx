import "./index.scoped.css";
function AppContainer(props) {
    return (
        <div className="main">
         <div className="content">
          {props.children}
        </div>
      </div>
    );
}  
export default AppContainer;