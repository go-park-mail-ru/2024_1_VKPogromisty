while IFS= read -r file; do
    if [[ -f "./static/$file" ]]; then
        echo "$file exists"
        ./mc cp "./static/$file" socio-minio/post-attachments
        rm "./static/$file"
    else
        echo "$file does not exist"
    fi
done < post-attachments.txt
