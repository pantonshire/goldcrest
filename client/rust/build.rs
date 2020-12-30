fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("../../protocol/twitter1.proto")?;
    Ok(())
}