Summarize this document:

Building Graph RAG for structured and unstructured data.
Sreeram A S K
Towards AI
Sreeram A S K

·
Follow

Published in
Towards AI

·
6 min read
·
Jan 11, 2025
48







RAG architecture is, by far, the most adapted and sophisticated solution for missing contextualisation of LLM’s. With no overhead of fine tuning, to a huge extent problems concerning the usage of LLM’s with untrained knowledge base’s have been solved with RAG.

While the Vector RAG could establish the contextualisation, the extent to which it could do that, has been limited. With complicated relationships and highly interconnected data, the recall measure of Vector RAG is not impressive. One of the major reasons being, the naïve vector embeddings that make up the Knowledge Base, which only consider geometrical proximity.

Graphs, on the other hand are innately structured to capture intricate relationships within data, leading to longer contextuality. For this reason graph based RAG’s have become the best means to exploit LLM capabilities.

Building Knowledge Graphs
Knowledge graphs can be populated with both unstructured data like Text, PDF’s and structured data like Tables/CSV’s. To manually identify and extract Nodes, Relationships and Node properties from documents spanning multiple pages, is a in-human task and requires high amount of domain expertise. Langchain has turned this daunting task into an effortless one, by using LLM’s for graph entity extraction.

Lets take a look at how to convert Unstructured and Structured data into Knowledge Graphs.

Unstructured data
Before proceeding further lets make all the necessary imports

from langchain_community.graphs import Neo4jGraph
from langchain_experimental.graph_transformers import LLMGraphTransformer
from langchain_openai import AzureChatOpenAI
import os
from langchain_core.documents import Document
from langchain.text_splitter import RecursiveCharacterTextSplitter, CharacterTextSplitter
from langchain_community.vectorstores.neo4j_vector import Neo4jVector
from langchain_openai import AzureOpenAIEmbeddings
import fitz
import logging 
import ast
from tqdm.notebook import tqdm
from concurrent.futures import ThreadPoolExecutor, as_completed
from langchain.memory import ConversationBufferMemory
We will be using Neo4j Desktop application as Graph Database and AzureOpenAI LLM across this article.

#Azure opena ai Gpt model initialisation
gpt_llm = AzureChatOpenAI(temperature=0, azure_deployment= OPENAI_GPT_DEPLOYMENT_NAME, api_key= OPENAI_API_KEY, api_version= OPENAI_API_VERSION, azure_endpoint= AZURE_ENDPOINT)
# Initialise neo4j object with an active neo4j server credentials 
graph_db = Neo4jGraph()
First we need to extract all the text from PDF’s and instantiate a langchain Document object with that text.

def convert_pdf_to_text (file_paths):

    """
    file_paths = list[str]
    rtype: Document object 
    Implements text extraction from pdf and forms a langchain object with the document text.
    """
  
    doc_text = []
    #iterate through each document
    for i in file_paths:
        logging.info("Reading ..."+str(i))
        doc = fitz.open(i)
            
        # Iterate through each page in the document and extract text
        try:
            text = ""
            for page_num in range(len(doc)):
                page = doc.load_page(page_num)
                text += page.get_text()
            # Combine metadata along with the document text  
            doc_text.append([Document(page_content=str(text))])
        
        except Exception as e:
            logging.info("Error Reading Document : "+str(e))

    return(doc_text)
After the text extraction, we now chunk the documents. Each new chunk will be a Document object in itself.

def split_text(docs):

    """
    docs: list of docuemnt objects 
    rtype: list[Document]
    Implements chunking mechanism by splitting the given document into chunks along with the metadata 
    """
    #initialise text splitter with appropriate chunk size and chunk overlap
    text_splitter = RecursiveCharacterTextSplitter.from_tiktoken_encoder(
    chunk_size=1000, chunk_overlap=100 )
    documents = []
    #Iterate through all the documents and implement chunking
    for i in docs:
        documents.append(text_splitter.split_documents(i))
    # merge all the chunked documents into a single list
    merged_docs = [item for sublist in documents for item in sublist]

    return(merged_docs)
The chunk and overlap size can be treated as hyper parameters.

We now have to populate a Graph from the chunked documents by extracting Nodes, Node properties and Relationships from the Document text. As the number of documents increase, it would become an impossible task to extract them manually. Hence we would use an LLM for extraction. The LLM would go through all the document chunks and populate those nodes which have interconnected relationships. The LLMGraphTransformer function from langchain does the job for us.

def construct_graph(doc):
    """ 
    doc: Document object
    rtype: Graph Document
    Implements graph contruction functionality using llm.
    """
    #initialise graph transformer object
    llm_transformer = LLMGraphTransformer(llm= gpt_llm)
    
    #construct graph by extracting nodes and relationships from documents using LLM
    graph_doc = llm_transformer.convert_to_graph_documents([doc])
    return(graph_doc)
Though this might have reduced a lot of manual labour, it still consumes a huge amount time for entity extraction. Hence we would multi thread the same process to handle multiple documents and collate the results. Below is a multi-threaded implementation of the above function.

def  thread_construct_graph(merged_docs):
  
    """
    merged_docs: list
    Implements the graph construction from diocuments using LLM in a multithreaded approach and pushes the extracted nodes and relationships into the provided GraphDB
    """
    MAX_WORKERS = 20 # can be changed as per requirement
    # NUM_ARTICLES = len(merged_docs)
    graph_documents = []

    with ThreadPoolExecutor(max_workers=MAX_WORKERS) as pool:
        # Submitting all tasks and creating a list of future objects
        futures =  [pool.submit(construct_graph, merged_docs[i]) for i in range(len(merged_docs)) ]
    
        #tqdm displays a progress bar for threading execution visualisation 
        for future in tqdm(as_completed(futures), total=len(futures), desc="Processing documents"):
            
            #capture results of thread status
            graph_document = future.result()
            graph_documents.extend(graph_document) #list extension method to add new graph documents to the existing list

    # add the contructed graph into graph db with extracted nodes and relatiosnhips 
    graph_db.add_graph_documents(graph_documents,baseEntityLabel = True, include_source = True)
    logging.info("Constructed graphdb...")
    logging.info("Schema: "+str(graph_db.get_schema))
Once the graph has been populated the add_graph_documents method of the graph_db object would push the populated graph into neo4j DB.

Structured Data
Constructing Knowledge Graph from structured data like CSV’s is comparatively strenuous than unstructured, but can be programmatically populated. The neomodel library, which connects to neo4j from python, can be used to define the structure of the node and its properties.

Required imports:

from neomodel import db, config, StructuredNode, RelationshipTo, RelationshipFrom, StringProperty
config.DATABASE_URL = "bolt://neo4j:#PASSWORD@localhost:7687"
To define the structure of Node, we will have to define a node class, implementing the StructuredNode. The member variables of the class will become the node’s properties, and relationships.

The below code demonstrates , an employee node that has properties such as Emp_ID, Emp_Fname, Emp_Lname, Emp_Location, Emp_Manager_Name, Emp_Rank..etc which is related to Manager Node by ‘reports to’ relation and to the Department Node by ‘belongs to’ relationship.

class Employee (StructuredNode):
    
    __label__ = "Employee"
    Emp_ID= StringProperty(unique_index = True)
    Emp_Fname= StringProperty()
    Emp_Lname= StringProperty()
    Emp_Location= StringProperty()
    Emp_Manager_Name= StringProperty()
    Emp_Rank= StringProperty()
    
    reports_to= RelationshipTo(Manager, "reports to")
    belongs_to= RelationshipTo(Department, "belongs to")


class Manager (StructuredNode):
    
    __label__ = "Manager"
    Emp_id= StringProperty(unique_index = True)
    M_Fname= StringProperty()
    M_Lname= StringProperty()
    dept_id=  StringProperty()
    dept_name = StringProperty()

class Department (StructuredNode):
  
    __label__ = "Department"
    Dept_ID= StringProperty(unique_index = True)
    Dept_name= StringProperty()
    Emp_count = StringProperty()
    Manager_emp_id = StringProperty()
After defining the Node properties and relationships its allowed to have, we can now instantiate the nodes and assign values to all the node properties. For example, if i have to define a new employee node and establish a relationship between their Manager and Department nodes, here is how we could do it.

emp_node = Employee.get_or_create({
    Emp_ID= #employee_id
    Emp_Fname= #employee_fname
    Emp_Lname= #employee_lname
    Emp_Location= #employee_location
    Emp_Manager_Name= .#employee_manager_name
    Emp_Rank= #employee_rank
    })

Manager = Manager.get_or_create({
    Emp_id= #manager_emp_id
    M_Fname= #manager_fname
    M_Lname= #manager_lname
    dept_id= #dept_id
    dept_name = #dept_name 
})

dept = Department.get_or_create({
    Dept_ID= #dept_id
    Dept_name= #dept_name
    Emp_count = #emp_count
    Manager_emp_id =#manager_emp_id
})


emp_node.reports_to.connect(Manager)
emp_node.belongs_to.connect(dept)

The get_or_create function either fetches an existing node if already present or creates a new node. If the manager or department nodes for an employee node already exist, then it assigns the employee node to them with the respective relationships.

LLM Response
Once the knowledge graph is developed, we would now be completing the Graph RAG cycle by developing the last module, that connects Graph DB with an LLM and answers user queries.

Langchain has several libraries which in connection with neo4j provides an interface to help build an QA chain. Lets explore one of those functions below which queries the Graph db with the user query

def query_db_with_llm(existing_graph_index,  metadata, query):

    """
    existing_graph_index = Neo4j Object
    metadata= dict
    query= str
    Implements querying the Graph DB using an LLM, alongside implicit metadata filtering
    """
  
    
    logging.info("Metadata = ", metadata)
    logging.info("Query = ", query)
    
    #Initialise Neo4j object from the existing kowledge graph
    existing_graph_index = Neo4jVector.from_existing_graph(
    az_embeddings, # use the same embedding model as used for the embedding creations
    )
    # Intialise langchain querying object
    qa_chain = ConversationalRetrievalChain.from_llm(return_source_documents= True, # True when the source of answers is required else False
    llm= #LLM of your choice, 
    retriever = existing_graph_index.as_retriever(search_kwargs={'filter': metadata}), #does a vector embedding semantic search along with metadata filtering on node properties using neo4j vector indexing 
    verbose = True,
    )
    #query the Graph db 
    answer = qa_chain({"question": query, "chat_history": []}) 
    
    return(answer)
The above code also demonstrates an additional metadata filtering functionality which filters nodes & relationships based on their properties, by passing it as a argument to index.as_retriever function.

Embedding
As many of you have noticed, we have not explicitly created embeddings for the unstructured text data before inserting it into Graph db. The LLMGraphTransformer would by default create a node label by name Document wherein the nodes under it would contain chunks and their embeddings as node properties. The embeddings are generated when the Neo4jVector.from_existing_graph method is called, if not already exists.

*In my next article i shall discuss and propose an enhanced Graph RAG querying method. Stay tuned!!!

Hope you find this informative.

AI